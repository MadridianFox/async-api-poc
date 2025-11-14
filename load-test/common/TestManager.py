from clickhouse_connect.driver import Client
from clickhouse_connect.driver.types import Matrix

class TestManager:
    def __init__(self, client: Client):
        self.client = client

    def list_tests(self) -> Matrix:
        return self.client.query('''
            SELECT test_name, test_id, timestamp 
            FROM load_tests 
            ORDER BY timestamp DESC, test_name ASC
        ''').result_rows

    def delete_test(self, test_id: str):
        self.client.command(f"DELETE from requests where test_id = '{test_id}'")
        self.client.command(f"DELETE from user_count where test_id = '{test_id}'")
        self.client.command(f"DELETE from load_tests where test_id = '{test_id}'")