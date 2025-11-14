import logging
import re
from datetime import datetime

from gevent import Greenlet, spawn, sleep
from clickhouse_connect.driver import Client
from locust import events
from locust.env import Environment

from common.LoadTest import LoadTest


class ClickhouseResultsPlugin:
    def __init__(self, load_test: LoadTest, client: Client, send_interval_s: int = 5, path_patterns: list = None):
        self.load_test = load_test
        self.client = client
        self.log = logging.getLogger('clickhouse_results_plugin')
        self.stop_flag = False
        self.send_interval_s = send_interval_s
        self.path_patterns = [] if path_patterns is None else path_patterns
        self.requests: list = []
        self.env: Environment

        self._create_tables()

        events.init.add_listener(lambda *args, **kwargs: self._init_handler(*args, **kwargs))
        events.request.add_listener(lambda *args, **kwargs: self._request_handler(*args, **kwargs))
        events.quitting.add_listener(lambda *args, **kwargs: self._quit_handler(*args, **kwargs))

        self.worker: Greenlet = spawn(self._run)

    def _run(self):
        self.log.info('Worker started.')
        while not self.stop_flag:
            sleep(self.send_interval_s)
            self._save_metrics()
            self._save_users_count()

    def _init_handler(self, environment, **_):
        self.env = environment

    def _request_handler(self, request_type, name, response_time, response_length, response,
                         context, exception, start_time, url, **kwargs):
        status = "success" if exception is None else "failure"
        error_message = "" if exception is None else repr(exception)
        request_path = name.split('?')[0]
        for pattern in self.path_patterns:
            request_path = re.sub(pattern[0], pattern[1], request_path)

        self.requests.append({
            "test_id": self.load_test.test_id,
            "tag": context["tag"] if "tag" in context else "no-tag",
            "timestamp": int(start_time),
            "endpoint": f"{request_type} {request_path}",
            "response_time": response_time,
            "status": status,
            "http_code": response.status_code,
            "exception": error_message
        })

    def _quit_handler(self, *_, **__):
        self.stop_flag = True
        self.worker.join()
        self._save_metrics()
        self._save_users_count()

    def _save_metrics(self):
        if len(self.requests) < 1:
            return

        data = [list(req.values()) for req in self.requests]
        fields = self.requests[0].keys()
        self.client.insert('requests', data, column_names=fields)

        self.requests.clear()

    def _save_users_count(self):
        if self.env.runner is None:
            return

        self.client.insert(
            'user_count',
            [
                [
                    self.load_test.test_id,
                    int(datetime.now().timestamp()),
                    self.env.runner.user_count
                ]
            ],
            ["test_id", "timestamp", "value"]
        )

    def _create_tables(self):
        self.client.command('''
            CREATE TABLE IF NOT EXISTS requests (
                test_id String,
                tag String,
                timestamp DateTime,
                endpoint String,
                response_time Float32,
                status String,
                http_code UInt16,
                exception String,
            ) ENGINE = MergeTree()
            ORDER BY timestamp;
        ''')

        self.client.command('''
            CREATE TABLE IF NOT EXISTS user_count (
                test_id String,
                timestamp DateTime,
                value UInt16,
            ) ENGINE = MergeTree()
            ORDER BY timestamp;
        ''')
