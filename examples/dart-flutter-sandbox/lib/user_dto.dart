/// Data Transfer Object representing a user profile model.
class UserDto {
  final String id;
  final String emailAddress;
  final String displayName;

  const UserDto({
    required this.id,
    required this.emailAddress,
    required this.displayName,
  });

  /// Parse user data from raw API JSON payload.
  factory UserDto.fromJson(Map<String, dynamic> json) {
    return UserDto(
      id: json['id'] as String? ?? '',
      emailAddress: json['email_address'] as String? ?? '',
      displayName: json['display_name'] as String? ?? '',
    );
  }

  /// Serialize model back to API JSON payload format.
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email_address': emailAddress,
      'display_name': displayName,
    };
  }
}
