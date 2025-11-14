<?php

namespace App\Requests;

use Hyperf\Validation\Request\FormRequest;

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
        return $this->input('user_id');
    }

    public function getProductId(): int
    {
        return $this->input('product_id');
    }

    public function getQty(): int
    {
        return $this->input('qty');
    }

    public function isWithItems(): bool
    {
        return $this->input('with_items', false);
    }
}
