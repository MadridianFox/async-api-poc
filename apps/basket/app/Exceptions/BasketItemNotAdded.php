<?php

namespace App\Exceptions;

use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class BasketItemNotAdded extends \Exception
{
    public function report(): void
    {}

    public function render(Request $request): JsonResponse
    {
        return response()->json([
            'code' => 'BasketItemNotAdded',
            'message' => 'Basket Item Not Added',
        ], 400);
    }
}
