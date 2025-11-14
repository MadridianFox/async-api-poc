<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

/**
 * @property int $id
 * @property int $basket_id
 * @property int $product_id
 * @property int $price_by_one
 * @property int $qty
 *
 * @property Basket $basket
 */
class BasketItem extends Model
{
    use HasFactory;

    protected $attributes = [
        'price_by_one' => 0,
        'qty' => 0,
    ];

    public function basket(): BelongsTo
    {
        return $this->belongsTo(Basket::class);
    }
}
