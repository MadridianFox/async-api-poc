<?php

declare(strict_types=1);

use function Hyperf\Support\env;

return [
    'basket' => [
        'host' => env("MS_BASKET_HOST", 'http://127.0.0.1/'),
    ],
    'product' => [
        'host' => env("MS_PRODUCT_HOST", 'http://127.0.0.1/'),
    ]
];
