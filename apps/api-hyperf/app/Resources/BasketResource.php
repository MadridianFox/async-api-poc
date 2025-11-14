<?php

namespace App\Resources;

use App\Domain\Basket\BasketDto;
use Hyperf\Resource\Json\JsonResource;

/** @mixin BasketDto */
class BasketResource extends JsonResource
{
    public function toArray(): array
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
