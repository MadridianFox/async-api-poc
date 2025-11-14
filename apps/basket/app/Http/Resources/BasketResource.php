<?php

namespace App\Http\Resources;

use App\Models\Basket;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/** @mixin Basket */
class BasketResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'user_id' => $this->user_id,
            'price' => $this->price,
            'status' => $this->status,

            'items' => $this->whenLoaded('items', fn () => BasketItemResource::collection($this->items))
        ];
    }
}
