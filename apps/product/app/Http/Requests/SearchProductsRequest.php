<?php

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;

class SearchProductsRequest extends FormRequest
{
    public function rules(): array
    {
        return [
            'price_min' => 'nullable|numeric|min:0',
            'price_max' => 'nullable|numeric|min:0',
            'name_like' => 'nullable|string|min:1|max:255',
            'qty_min' => 'nullable|numeric|min:0',
            'qty_max' => 'nullable|numeric|min:0',
        ];
    }

    public function getPriceMin(): ?int
    {
        return $this->input('price_min');
    }

    public function getPriceMax(): ?int
    {
        return $this->input('price_max');
    }

    public function getNameLike(): ?string
    {
        return $this->input('name_like');
    }

    public function getQtyMin(): ?int
    {
        return $this->input('qty_min');
    }

    public function getQtyMax(): ?int
    {
        return $this->input('qty_max');
    }
}
