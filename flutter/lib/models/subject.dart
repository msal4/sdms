import 'dart:convert';

import 'package:sdms/models/lecturer.dart';

Subject subjectFromJson(String str) => Subject.fromJson(json.decode(str));

String subjectToJson(Subject data) => json.encode(data.toJson());

class Subject {
  Subject({
    required this.id,
    required this.name,
    this.details,
    required this.stage,
    required this.lecturer,
    required this.semester,
    this.syllabus,
  });

  int id;
  String name;
  String? details;
  int stage;
  Lecturer lecturer;
  int semester;
  String? syllabus;

  factory Subject.fromJson(Map<String, dynamic> json) => Subject(
        id: json["ID"],
        name: json["Name"],
        details: json["Details"],
        stage: json["Stage"],
        lecturer: Lecturer.fromJson(json["Lecturer"]),
        semester: json["Semester"],
        syllabus: json["Syllabus"],
      );

  Map<String, dynamic> toJson() => {
        "ID": id,
        "Name": name,
        "Details": details,
        "Stage": stage,
        "Lecturer": lecturer.toJson(),
        "Semester": semester,
        "Syllabus": syllabus,
      };
}
