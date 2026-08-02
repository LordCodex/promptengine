<?php

declare(strict_types=1);

namespace App\Http\Controllers;

use App\Actions\RegisterUserAction;
use App\Http\Requests\StoreUserRequest;
use Illuminate\Http\JsonResponse;

class UserController
{
    public function __construct(
        private RegisterUserAction $registerUser
    ) {}

    /**
     * Store a newly created user profile.
     */
    public function store(StoreUserRequest $request): JsonResponse
    {
        $user = $this->registerUser->execute(
            $request->validated()
        );

        return response()->json([
            'data' => [
                'id' => $user->id,
                'email' => $user->email,
                'name' => $user->name,
            ]
        ], 201);
    }
}
