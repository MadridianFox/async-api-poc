<?php

namespace App\Domains\Catalog;

use GuzzleHttp\Client;

class ProductClient
{
    private readonly Client $client;
    public function __construct()
    {
        $this->client = new Client([
            'base_uri' => config('ms.product.host'),
            'http_errors' => false,
        ]);
    }

    public function searchProducts(ProductFilter $filter): ProductSearchResult
    {
        $response = $this->client->get("/api/products", ["query" => [
            'price_min' => $filter->getMinPrice(),
            'price_max' => $filter->getMaxPrice(),
            'name_like' => $filter->getNameLike(),
            'qty_min' => $filter->getMinQty(),
            'qty_max' => $filter->getMaxQty(),
            'page' => $filter->getPage(),
        ]]);

        if ($response->getStatusCode() !== 200) {
            return new ProductSearchResult([], []);
        }
        $body = json_decode($response->getBody()->getContents(), true);

        return new ProductSearchResult($body['data'] ?? [], $body['meta'] ?? []);
    }
}
