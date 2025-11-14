<?php

namespace App\Domain\Basket;

use GuzzleHttp\Client;
use Hyperf\Contract\ConfigInterface;
use Hyperf\Guzzle\ClientFactory;
use Hyperf\Logger\LoggerFactory;
use Psr\Log\LoggerInterface;

class BasketClient
{
    private readonly Client $client;
    private LoggerInterface $logger;

    public function __construct(
        ClientFactory $clientFactory,
        LoggerFactory $loggerFactory,
        ConfigInterface $config
    )
    {
        $this->logger = $loggerFactory->get('log');

        $this->client = $clientFactory->create([
            'base_uri' => $config->get('ms.basket.host'),
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
            $this->logger->error("Basket bad response for /baskets/current", [
                "status" => $response->getStatusCode(),
                "body" => $bodyStr
            ]);

            return null;
        }

        return new BasketDto($body['data']);
    }
}