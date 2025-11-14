<?php

namespace App\Requests;

use App\Domain\Catalog\ProductFilter;
use Hyperf\Validation\Request\FormRequest;

class SearchProductsRequest extends FormRequest implements ProductFilter
{

    public function rules(): array
    {
        return [
            'price_min' => 'nullable|numeric|min:0',
            'price_max' => 'nullable|numeric|min:0',
            'name_like' => 'nullable|string|min:1|max:255',
            'qty_min' => 'nullable|numeric|min:0',
            'qty_max' => 'nullable|numeric|min:0',
            'page' => 'nullable|numeric|min:1',
        ];
    }
    public function getMinQty(): ?int
    {
        return $this->input('qty_min');
    }

    public function getMaxQty(): ?int
    {
        return $this->input('qty_max');
    }

    public function getMinPrice(): ?int
    {
        return $this->input('price_min');
    }

    public function getMaxPrice(): ?int
    {
        return $this->input('price_max');
    }

    public function getNameLike(): ?string
    {
        return $this->input('name_like');
    }

    public function getPage(): ?int
    {
        return $this->input('page');
    }
}
