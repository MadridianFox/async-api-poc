<?php

namespace App\Http\Resources;

use App\Models\BasketItem;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/** @mixin BasketItem */
class BasketItemResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'product_id' => $this->product_id,
            'price_by_one' => $this->price_by_one,
            'qty' => $this->qty,
        ];
    }
}
