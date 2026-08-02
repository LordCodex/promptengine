<?php

declare(strict_types=1);

namespace App\Actions;

use App\Models\User;
use Illuminate\Support\Facades\DB;
use InvalidArgumentException;

class RegisterUserAction
{
    /**
     * Execute the register user transaction.
     *
     * @param array{email: string, name: string} $data
     * @throws InvalidArgumentException
     */
    public function execute(array $data): User
    {
        if (empty($data['email'])) {
            throw new InvalidArgumentException('Email field cannot be empty.');
        }

        return DB::transaction(function () use ($data) {
            return User::create([
                'email' => $data['email'],
                'name' => $data['name'],
                'is_active' => true,
            ]);
        });
    }
}
