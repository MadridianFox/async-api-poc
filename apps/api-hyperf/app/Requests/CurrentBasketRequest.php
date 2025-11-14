<?php

namespace App\Requests;


use Hyperf\Validation\Request\FormRequest;

class CurrentBasketRequest extends FormRequest
{
    public function rules(): array
    {
        return [
            'user_id' => 'required|integer',
            'with_items' => 'nullable|boolean',
        ];
    }

    public function getUserId(): int
    {
        return $this->query('user_id');
    }

    public function isWithItems(): bool
    {
        return $this->query('with_items', false);
    }
}
