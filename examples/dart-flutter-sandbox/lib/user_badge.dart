import 'package:flutter/material.dart';

/// Rebuild-optimized widget to display a user's name badge.
class UserBadge extends StatelessWidget {
  final String displayName;

  const UserBadge({
    super.key,
    required this.displayName,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 12.0, vertical: 6.0),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Icon(Icons.person, size: 16.0),
          const SizedBox(width: 8.0),
          Text(
            displayName,
            style: const TextStyle(
              fontSize: 14.0,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }
}
