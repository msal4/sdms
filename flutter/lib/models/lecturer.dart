// To parse this JSON data, do
//
//     final lecturer = lecturerFromJson(jsonString);

import 'dart:convert';

Lecturer lecturerFromJson(String str) => Lecturer.fromJson(json.decode(str));

String lecturerToJson(Lecturer data) => json.encode(data.toJson());

class Lecturer {
  Lecturer({
    required this.id,
    required this.name,
    this.image,
    required this.username,
    required this.password,
    this.about,
  });

  int id;
  String name;
  String? image;
  String username;
  String password;
  String? about;

  factory Lecturer.fromJson(Map<String, dynamic> json) => Lecturer(
        id: json["ID"],
        name: json["Name"],
        image: json["Image"],
        username: json["Username"],
        password: json["Password"],
        about: json["About"],
      );

  Map<String, dynamic> toJson() => {
        "ID": id,
        "Name": name,
        "Image": image,
        "Username": username,
        "Password": password,
        "About": about,
      };
}
