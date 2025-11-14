<?php

namespace App\Resources;

use App\Domain\Catalog\ProductDto;
use Hyperf\Resource\Json\JsonResource;

/** @mixin ProductDto */
class ProductResource extends JsonResource
{
    public function toArray(): array
    {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'price' => $this->price,
            'qty' => $this->qty,
            'created_at' => $this->created_at->toIso8601String(),
            'updated_at' => $this->updated_at->toIso8601String(),
        ];
    }
}
