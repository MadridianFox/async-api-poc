<?php

namespace App\Http\Controllers;

use App\Actions\SetItemAction;
use App\Http\Requests\CurrentBasketRequest;
use App\Http\Requests\SetItemRequest;
use App\Http\Resources\BasketResource;
use App\Models\Basket;
use App\Models\BasketItem;
use Illuminate\Http\Resources\Json\JsonResource;
use Illuminate\Support\Facades\DB;

class BasketController extends Controller
{
    public function getCurrentBasket(CurrentBasketRequest $request): JsonResource
    {
        // simulate hard work
        usleep(rand(50_000, 400_000));

        $basket = $this->findCurrentBasketOrNew($request->getUserId(), $request->isWithItems());

        return BasketResource::make($basket);
    }

    public function setItem(SetItemRequest $request, SetItemAction $action): JsonResource
    {
        // simulate hard work
        usleep(rand(500_000, 3000_000));

        $action->execute(
            $this->findCurrentBasketOrNew($request->getUserId(), true),
            $request->getProductId(),
            $request->getQty()
        );

        $basket = $this->findCurrentBasketOrNew($request->getUserId(), $request->isWithItems());
        return BasketResource::make($basket);
    }

    protected function findCurrentBasketOrNew(int $userId, bool $withItems): Basket
    {
        $query = Basket::query()
            ->where('user_id', $userId)
            ->where('status', Basket::STATUS_NEW)
            ->orderBy('created_at', 'desc');

        if ($withItems) {
            $query->with('items');
        }

        /** @var Basket $basket */
        $basket = $query->first() ?? new Basket();
        $basket->user_id = $userId;

        return $basket;
    }
}
