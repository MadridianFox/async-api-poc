<?php

namespace App\Domain\Basket;

class BasketItemDto
{
    public readonly int $id;
    public readonly int $productId;
    public readonly int $priceByOne;
    public readonly int $qty;

    public function __construct(array $data)
    {
        $this->id = $data['id'];
        $this->productId = $data['product_id'];
        $this->priceByOne = $data['price_by_one'];
        $this->qty = $data['qty'];
    }
}
