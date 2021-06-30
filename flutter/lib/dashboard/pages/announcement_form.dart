import 'dart:typed_data';

import 'package:dio/dio.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:http_parser/http_parser.dart';
import 'package:sdms/client.dart';
import 'package:sdms/const.dart';
import 'package:sdms/models/announcement.dart';

class AnnouncementFormPage extends StatefulWidget {
  const AnnouncementFormPage({Key? key, this.next}) : super(key: key);

  String get title => "Add Announcement";

  final VoidCallback? next;

  @override
  _AnnouncementFormPageState createState() => _AnnouncementFormPageState();
}

class _AnnouncementFormPageState extends State<AnnouncementFormPage> {
  late final TextEditingController _titleController;
  late final TextEditingController _detailsController;
  Uint8List? _imageBytes;

  @override
  void initState() {
    _titleController = TextEditingController();
    _detailsController = TextEditingController();

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
              controller: _titleController,
              decoration:
                  InputDecoration(labelText: "Title of the Announcement"),
            ),
            SizedBox(height: 10),
            TextField(
              controller: _detailsController,
              decoration: InputDecoration(labelText: "Details"),
            ),
            SizedBox(height: 10),
            Row(
              crossAxisAlignment: CrossAxisAlignment.center,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                FlatButton.icon(
                    onPressed: () async {
                      final result = await FilePicker.platform.pickFiles(
                          type: FileType.custom,
                          allowedExtensions: ["png", "jpg", "jpeg"]);
                      if (result == null || result.count == 0) return;
                      setState(() {
                        _imageBytes = result.files.single.bytes;
                      });
                      debugPrint("File is attached");
                    },
                    icon: Icon(Icons.upload),
                    label: Text("Upload Image")),
                IconButton(
                  color: Colors.grey,
                  icon: Icon(Icons.delete,
                      color: _imageBytes != null ? Colors.red : null),
                  onPressed: () {
                    setState(() {
                      _imageBytes = null;
                    });
                  },
                ),
              ],
            ),
            SizedBox(height: 10),
            ElevatedButton.icon(
              onPressed: () {
                final an = Announcement(
                  title: _titleController.text,
                  details: _detailsController.text,
                ).toJson();

                if (_imageBytes != null) {
                  an["Image"] = MultipartFile.fromBytes(_imageBytes!,
                      filename: "image.png",
                      contentType: MediaType.parse("image/png"));
                }

                final data = FormData.fromMap(an);

                client
                    .post("/announcements",
                        data: data,
                        options: Options(
                            headers: {"content-type": "multipart/form-data"}))
                    .then((res) {
                  if (widget.next != null) widget.next!();
                  Navigator.pop(context);
                });
              },
              icon: Icon(Icons.add),
              label: Text("Add Announcement"),
            ),
          ],
        ),
      ),
    );
  }
}
