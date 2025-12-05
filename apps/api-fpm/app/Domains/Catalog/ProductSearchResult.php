<?php

namespace App\Domains\Catalog;

use Illuminate\Support\Collection;

class ProductSearchResult
{
    public readonly Collection $products;
    public readonly int $total;
    public readonly int $perPage;
    public readonly int $currentPage;

    public function __construct(array $data, array $meta)
    {
        $this->products = collect($data)->map(fn (array $product) => new ProductDto($product));
        $this->total = $meta['total'] ?? 0;
        $this->perPage = $meta['per_page'] ?? 0;
        $this->currentPage = $meta['current_page'] ?? 0;
    }
}
