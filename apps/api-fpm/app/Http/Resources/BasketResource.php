<?php

namespace App\Http\Resources;

use App\Domains\Basket\BasketDto;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/** @mixin BasketDto */
class BasketResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'user_id' => $this->userId,
            'price' => $this->price,
            'status' => $this->status,

            'items' => BasketItemResource::collection($this->items),
        ];
    }
}
