<?php

namespace App\Domain\Catalog;

use GuzzleHttp\Client;
use Hyperf\Contract\ConfigInterface;
use Hyperf\Guzzle\ClientFactory;

class CatalogClient
{
    private readonly Client $client;
    public function __construct(ClientFactory $clientFactory, ConfigInterface $config)
    {
        $this->client = $clientFactory->create([
            'base_uri' => $config->get('ms.product.host'),
            'http_errors' => false,
            'headers' => [
                'Accept' => 'application/json',
            ]
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