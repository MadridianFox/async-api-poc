<?php

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;

class SetBasketItemRequest extends FormRequest
{
    public function rules(): array
    {
        return [
            'user_id' => 'required|integer',
            'product_id' => 'required|integer',
            'qty' => 'required|integer',
            'with_items' => 'nullable|boolean',
        ];
    }

    public function getUserId(): int
    {
        return $this->get('user_id');
    }

    public function getProductId(): int
    {
        return $this->get('product_id');
    }

    public function getQty(): int
    {
        return $this->get('qty');
    }

    public function isWithItems(): bool
    {
        return $this->get('with_items', false);
    }
}
