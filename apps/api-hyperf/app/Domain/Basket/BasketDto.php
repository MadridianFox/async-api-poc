<?php

namespace App\Domain\Basket;

class BasketDto
{
    public readonly ?int $id;
    public readonly int $userId;
    public readonly int $price;
    public readonly int $status;
    public readonly array $items;

    public function __construct(array $data)
    {
        $this->id = $data['id'];
        $this->userId = $data['user_id'];
        $this->price = $data['price'];
        $this->status = $data['status'];
        $this->items = $this->parseItems($data['items'] ?? []);
    }

    private function parseItems(array $rawItems): array
    {
        return array_map(function ($item) {
            return new BasketItemDto($item);
        }, $rawItems);
    }
}
