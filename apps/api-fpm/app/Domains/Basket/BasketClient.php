<?php

namespace App\Domains\Basket;

use App\Domains\Catalog\ProductSearchResult;
use GuzzleHttp\Client;

class BasketClient
{
    private readonly Client $client;
    public function __construct()
    {
        $this->client = new Client([
            'base_uri' => config('ms.basket.host'),
            'http_errors' => false,
            'headers' => [
                'Accept' => 'application/json',
            ]
        ]);
    }

    public function setItem(int $userId, int $productId, int $quantity, bool $withItems = false): ?BasketDto
    {
        $response = $this->client->post('/api/baskets/current/set-item', [
            'json' => [
                'user_id' => $userId,
                'product_id' => $productId,
                'qty' => $quantity,
                'with_items' => $withItems,
            ]
        ]);

        if ($response->getStatusCode() !== 200) {
            return null;
        }

        $body = json_decode($response->getBody()->getContents(), true);

        return new BasketDto($body['data']);
    }

    public function getCurrentBasket(int $userId, bool $withItems = false): ?BasketDto
    {
        $response = $this->client->get('/api/baskets/current', [
            'query' => [
                'user_id' => $userId,
                'with_items' => $withItems,
            ]
        ]);

        $bodyStr = $response->getBody()->getContents();
        $body = json_decode($bodyStr, true);

        if ($response->getStatusCode() !== 200 || !$body || !isset($body['data'])) {
            logger()->error("Basket bad response for /baskets/current", [
                "status" => $response->getStatusCode(),
                "body" => $bodyStr
            ]);

            return null;
        }

        return new BasketDto($body['data']);
    }
}
