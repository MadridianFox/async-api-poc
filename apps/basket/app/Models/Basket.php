<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Collection;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Support\Carbon;

/**
 * @property int $id
 * @property int $user_id
 * @property int $price
 * @property int $status
 *
 * @property Carbon $created_at
 * @property Carbon $updated_at
 *
 * @property Collection|BasketItem[] $items
 */
class Basket extends Model
{
    const STATUS_NEW = 0;
    const STATUS_PROCESSING = 1;
    const STATUS_DONE = 2;

    use HasFactory;

    protected $attributes = [
        'status' => self::STATUS_NEW,
        'price' => 0,
    ];

    public function items(): HasMany
    {
        return $this->hasMany(BasketItem::class);
    }
}
