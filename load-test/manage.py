#!/bin/env python3

import argparse
from yaml import safe_load
import clickhouse_connect
from common.TestManager import TestManager

def main():
    parser = argparse.ArgumentParser(
        description="Утилита для управления таблицами и тестами."
    )

    subparsers = parser.add_subparsers(dest="command", required=True)

    subparsers.add_parser("list_tests")
    delete_parser = subparsers.add_parser("delete_test", help="Delete a test by the test_id")
    delete_parser.add_argument("test_id")

    args = parser.parse_args()

    settings: dict

    with open('settings.yaml', 'r') as settings_file:
        settings = safe_load(settings_file)

    ch_client = clickhouse_connect.create_client(
        host=settings["clickhouse"]["host"],
        port=settings["clickhouse"]["port"],
        username=settings["clickhouse"]["username"],
        password=settings["clickhouse"]["password"],
        database=settings["clickhouse"]["database"]
    )

    test_manager = TestManager(ch_client)

    match args.command:
        case "delete_test":
            test_manager.delete_test(args.test_id)
        case "list_tests":
            rows = test_manager.list_tests()
            print(f"{'ID':<40} {'Name':<30} Date")
            for row in rows:
                print(f"{row[1]:<40} {row[0]:<30} {str(row[2])}")

if __name__ == "__main__":
    main()