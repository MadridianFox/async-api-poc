<?php

namespace App\Actions;

use App\Exceptions\BasketItemNotAdded;
use App\Models\Basket;
use App\Models\BasketItem;
use GuzzleHttp\Client;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\DB;

class SetItemAction
{
    public function execute(Basket $basket, int $productId, int $qty): void
    {
        $items = $basket->items?->keyBy('product_id') ?? collect();
        $currentItem = $items[$productId] ?? null;

        if (!$currentItem) {
            $price = $this->getProductPrice($productId);
            if ($price === null) {
                throw new BasketItemNotAdded();
            }

            $currentItem = new BasketItem();
            $currentItem->product_id = $productId;
            $currentItem->price_by_one = $price;
            $items[$productId] = $currentItem;
        }

        $currentItem->qty = $qty;

        $this->recalculateBasket($basket, $items);

        DB::transaction(function () use ($basket, $currentItem) {
            $basket->save();

            if ($currentItem->qty > 0) {
                $currentItem->basket_id = $basket->id;
                $currentItem->save();
            } else {
                if ($currentItem->exists) {
                    $currentItem->delete();
                }
            }
        });
    }

    private function getProductPrice(int $productId): ?int
    {
        $client = new Client(['base_uri' => config('ms.product.host'), 'http_errors' => false]);
        $response = $client->get("/api/products/{$productId}");
        if ($response->getStatusCode() !== 200) {
            return null;
        }
        $body = json_decode($response->getBody()->getContents(), true);
        $product = $body['data'] ?? [];

        return $product['price'] ?? null;
    }

    /**
     * @param Basket $basket
     * @param Collection<BasketItem> $items
     */
    private function recalculateBasket(Basket $basket, Collection $items): void
    {
        $totalPrice = 0;
        foreach ($items as $item) {
            $totalPrice += $item->price_by_one * $item->qty;
        }

        $basket->price = $totalPrice;
    }
}
