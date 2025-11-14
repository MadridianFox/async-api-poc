<?php

namespace App\Resources;

use App\Domain\Basket\BasketItemDto;
use Hyperf\Resource\Json\JsonResource;

/** @mixin BasketItemDto */
class BasketItemResource extends JsonResource
{
    public function toArray(): array
    {
        return [
            'id' => $this->id,
            'product_id' => $this->productId,
            'price_by_one' => $this->priceByOne,
            'qty' => $this->qty,
        ];
    }
}
