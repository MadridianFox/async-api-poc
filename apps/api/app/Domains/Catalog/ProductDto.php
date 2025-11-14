<?php

namespace App\Domains\Catalog;

use Illuminate\Support\Carbon;

class ProductDto
{
    public int $id;
    public string $name;
    public int $price;
    public int $qty;
    public Carbon $created_at;
    public Carbon $updated_at;
    public function __construct(array $item)
    {
        $this->id = $item['id'];
        $this->name = $item['name'];
        $this->price = $item['price'];
        $this->qty = $item['qty'];
        $this->created_at = Carbon::parse($item['created_at']);
        $this->updated_at = Carbon::parse($item['updated_at']);
    }
}
