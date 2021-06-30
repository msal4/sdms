import 'dart:typed_data';

import 'package:dio/dio.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:http_parser/http_parser.dart';
import 'package:sdms/client.dart';
import 'package:sdms/const.dart';
import 'package:sdms/models/lecturer.dart';
import 'package:sdms/models/subject.dart';

class SubjectFormPage extends StatefulWidget {
  const SubjectFormPage({Key? key, this.next}) : super(key: key);

  String get title => "Add Subject";

  final VoidCallback? next;

  @override
  _SubjectFormPageState createState() => _SubjectFormPageState();
}

class _SubjectFormPageState extends State<SubjectFormPage> {
  late final TextEditingController _nameController;
  late final TextEditingController _detailsController;
  late final TextEditingController _semesterController;
  late final TextEditingController _stageController;
  Uint8List? _syllabusBytes;

  List<Lecturer> _lecturers = [];
  Lecturer? _currentLecturer;

  @override
  void initState() {
    _nameController = TextEditingController();
    _detailsController = TextEditingController();
    _semesterController = TextEditingController();
    _stageController = TextEditingController();

    getLecturers().then((value) => setState(() {
          _lecturers = value;
        }));

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(widget.title)),
      body: Container(
        padding: const EdgeInsets.all(kDefaultPadding),
        child: ListView(
          children: [
            TextField(
              controller: _nameController,
              decoration: InputDecoration(labelText: "Name of the Subject"),
            ),
            SizedBox(height: 10),
            TextField(
              controller: _detailsController,
              decoration: InputDecoration(labelText: "Details"),
            ),
            SizedBox(height: 10),
            TextField(
              controller: _semesterController,
              decoration: InputDecoration(labelText: "Semester"),
            ),
            SizedBox(height: 10),
            TextField(
              controller: _stageController,
              decoration: InputDecoration(labelText: "Stage"),
            ),
            SizedBox(height: 10),
            DropdownButton(
              value: _currentLecturer?.id,
              hint: Text("Choose Lecturer"),
              items: [
                for (final l in _lecturers)
                  DropdownMenuItem(value: l.id, child: Text(l.name))
              ],
              onChanged: (id) => setState(() {
                _currentLecturer = _lecturers.firstWhere((i) => i.id == id);
              }),
            ),
            SizedBox(height: 10),
            Row(
              crossAxisAlignment: CrossAxisAlignment.center,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                FlatButton.icon(
                    onPressed: () async {
                      final result = await FilePicker.platform.pickFiles(
                          type: FileType.custom, allowedExtensions: ["pdf"]);
                      if (result == null || result.count == 0) return;
                      setState(() {
                        _syllabusBytes = result.files.single.bytes;
                      });
                      debugPrint("File is attached");
                    },
                    icon: Icon(Icons.upload),
                    label: Text("Upload Syllabus")),
                IconButton(
                  color: Colors.grey,
                  icon: Icon(Icons.delete,
                      color: _syllabusBytes != null ? Colors.red : null),
                  onPressed: () {
                    setState(() {
                      _syllabusBytes = null;
                    });
                  },
                ),
              ],
            ),
            SizedBox(height: 10),
            ElevatedButton.icon(
              onPressed: () {
                if (_currentLecturer == null) return;

                final sub = Subject(
                  name: _nameController.text,
                  details: _detailsController.text,
                  semester: int.parse(_semesterController.text),
                  stage: int.parse(_stageController.text),
                ).toJson();

                sub["Lecturer"] = _currentLecturer!.id;

                if (_syllabusBytes != null) {
                  print("syllabus is not null");
                  sub["Syllabus"] = MultipartFile.fromBytes(
                    _syllabusBytes!,
                    filename: "testicles.pdf",
                    contentType: MediaType.parse("application/pdf"),
                  );
                }

                final data = FormData.fromMap(sub);

                client
                    .post("/subjects",
                        data: data,
                        options: Options(
                            headers: {"content-type": "multipart/form-data"}))
                    .then((res) {
                  if (widget.next != null) widget.next!();
                  Navigator.pop(context);
                });
              },
              icon: Icon(Icons.add),
              label: Text("Add Subject"),
            ),
          ],
        ),
      ),
    );
  }
}
