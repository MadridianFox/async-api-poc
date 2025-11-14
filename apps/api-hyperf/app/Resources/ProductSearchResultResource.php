<?php

namespace App\Resources;

use Hyperf\Resource\Json\JsonResource;

class ProductSearchResultResource extends JsonResource
{
    public function toArray(): array
    {
        return [
            'data' => ProductResource::collection($this->resource->products),
            'meta' => [
                'total' => $this->resource->total,
                'per_page' => $this->resource->perPage,
                'current_page' => $this->resource->currentPage,
            ]
        ];
    }
}
