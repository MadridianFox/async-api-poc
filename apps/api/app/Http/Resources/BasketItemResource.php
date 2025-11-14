<?php

namespace App\Http\Resources;

use App\Domains\Basket\BasketItemDto;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/** @mixin BasketItemDto */
class BasketItemResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'product_id' => $this->productId,
            'price_by_one' => $this->priceByOne,
            'qty' => $this->qty,
        ];
    }
}
