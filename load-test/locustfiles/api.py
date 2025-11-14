import random

from gevent import time
from yaml import safe_load
import clickhouse_connect

from locust import HttpUser, SequentialTaskSet, task, constant
from common.LoadTest import LoadTest
from common.ClickhouseResultsPlugin import ClickhouseResultsPlugin
from common.StepLoadShape import StepLoadShape

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

# Если в пути эндпоинта есть параметр, то это будет считаться разными эндпоинтами
# Чтобы этого избежать можно добавить в список регулярку чтобы схлопнуть все эти эндпоинты в один
# Первый элемент - регулярка, второй - то на что будет заменён путь
path_patterns = [
    # ["/api/v1/catalog/brands/.*", "/api/v1/catalog/brands/{id}"],
    # ["/api/v1/catalog/offers/.*", "/api/v1/catalog/offers/{id}"],
    # ["/api/v1/sections/products/.*", "/api/v1/sections/products/{id}"],
]


class LoadShape(StepLoadShape):
    def __init__(self):
        # Тут задаётся сколько параллельных пользователей в какой промежуток времени должно быть
        # первый элемент кортежа - длительность в секундах, второй - количество пользователей
        self.stages = [
            (1, 1),       # начинаем с одного пользователя, этот шаг лучше не удалять
            (30, 1),      # 30 секунд один пользователь
            (60, 100),    # в течение 60 секунд линейно поднимаем до 100 пользователей
            (120, 100),   # 120 секунд держим 100 пользователей
        ]
        super().__init__()

load_test = LoadTest("api", settings, ch_client)
clickhouse_results = ClickhouseResultsPlugin(load_test, ch_client, path_patterns=path_patterns)


class UnautorizedUser(HttpUser):
    wait_time = constant(1)

    @task(10)
    def get_products(self):
        request_body = {
            "qty_min": random.choice([
                None,
                random.randint(0, 100)
            ]),
            "qty_max": random.choice([
                None,
                random.randint(100, 1000)
            ]),
            "price_min": random.choice([
                None,
                random.randint(0, 1000)
            ]),
            "price_max": random.choice([
                None,
                random.randint(1000, 100000)
            ]),
            "page": random.randint(1,20)
        }
        self.client.post('api/catalog/search', json=request_body, headers={"Accept": "application/json"})

    @task(3)
    def get_basket(self):
        params = {
            "user_id": random.randint(1, 1000),
            "with_items": random.choice([0, 1])
        }
        self.client.get('api/basket/current', params=params, headers={"Accept": "application/json"})

    @task(1)
    def set_basket_item(self):
        params = {
            "user_id": random.randint(1, 1000),
            "with_items": random.choice([0, 1]),
            "product_id": random.randint(1, 10000),
            "qty": random.randint(1, 100),
        }
        self.client.post('api/basket/current/set-item', params=params, headers={"Accept": "application/json"})
