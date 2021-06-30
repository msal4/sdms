// To parse this JSON data, do
//
//     final announcement = announcementFromJson(jsonString);

import 'dart:convert';

Announcement announcementFromJson(String str) =>
    Announcement.fromJson(json.decode(str));

String announcementToJson(Announcement data) => json.encode(data.toJson());

class Announcement {
  Announcement({
    this.id,
    required this.title,
    this.details,
  });

  int? id;
  String title;
  String? details;

  factory Announcement.fromJson(Map<String, dynamic> json) => Announcement(
        id: json["ID"],
        title: json["Title"],
        details: json["Details"],
      );

  Map<String, dynamic> toJson() => {
        "ID": id,
        "Title": title,
        "Details": details,
      };
}
