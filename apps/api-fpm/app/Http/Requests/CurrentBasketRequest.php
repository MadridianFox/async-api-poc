<?php

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;

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
        return $this->get('user_id');
    }

    public function isWithItems(): bool
    {
        return $this->get('with_items', false);
    }
}
