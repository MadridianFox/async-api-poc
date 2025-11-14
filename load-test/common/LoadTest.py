import uuid
from datetime import datetime

from locust.env import Environment
from locust import events
from clickhouse_connect.driver import Client

class LoadTest:
    def __init__(self, name: str, settings:dict,  client: Client):
        self.name: str = name
        self.settings: dict = settings
        self.test_id: str = str(uuid.uuid4())
        self.client: Client = client
        self._create_tables()

        events.init.add_listener(lambda *args, **kwargs: self._init_handler(*args, **kwargs))

    def _init_handler(self, environment: Environment, *_, **__):
        self.env = environment

        self.client.insert("load_tests", [[
            int(datetime.now().timestamp()),
            self.name,
            self.test_id,
            self.env.host
        ]], ["timestamp", "test_name", "test_id", "url"])

    def _create_tables(self):
        self.client.command('''
            CREATE TABLE IF NOT EXISTS load_tests (
                timestamp DateTime,
                test_name String,
                test_id String,
                url String
            ) ENGINE = MergeTree()
            ORDER BY timestamp;
        ''')