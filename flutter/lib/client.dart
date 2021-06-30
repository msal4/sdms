import 'package:dio/dio.dart';
import 'package:sdms/models/announcement.dart';
import 'package:sdms/models/lecturer.dart';
import 'package:sdms/models/subject.dart';

final client = Dio(BaseOptions(baseUrl: "http://localhost:5000/api/"));

Future<List<Lecturer>> getLecturers() => client.get("/lecturers").then(
    (value) => (value.data as List).map((i) => Lecturer.fromJson(i)).toList());

Future<List<Subject>> getSubjects() => client.get("/subjects").then(
    (value) => (value.data as List).map((i) => Subject.fromJson(i)).toList());

Future<List<Announcement>> getAnnouncements() =>
    client.get("/announcements").then((value) =>
        (value.data as List).map((i) => Announcement.fromJson(i)).toList());
