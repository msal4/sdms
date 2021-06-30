import 'package:flutter/material.dart';
import 'package:sdms/client.dart';
import 'package:sdms/const.dart';
import 'package:sdms/models/lecturer.dart';

class LecturerFormPage extends StatefulWidget {
  const LecturerFormPage({Key? key, this.lecturer, this.next})
      : super(key: key);

  String get title => this.lecturer != null ? "Edit Lecturer" : "Add Lecturer";

  final Lecturer? lecturer;
  final VoidCallback? next;

  @override
  _LecturerFormPageState createState() => _LecturerFormPageState();
}

class _LecturerFormPageState extends State<LecturerFormPage> {
  late final TextEditingController _nameController;
  late final TextEditingController _usernameController;
  late final TextEditingController _passwordController;
  late final TextEditingController _aboutController;

  Lecturer? get lecturer => widget.lecturer;

  @override
  void initState() {
    _nameController = TextEditingController.fromValue(
        TextEditingValue(text: lecturer?.name ?? ""));
    _usernameController = TextEditingController.fromValue(
        TextEditingValue(text: lecturer?.username ?? ""));
    _passwordController = TextEditingController.fromValue(
        TextEditingValue(text: lecturer?.password ?? ""));
    _aboutController = TextEditingController.fromValue(
        TextEditingValue(text: lecturer?.about ?? ""));

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
              decoration: InputDecoration(labelText: "Name of the Lecturer"),
            ),
            SizedBox(height: 10),
            TextField(
              controller: _usernameController,
              decoration: InputDecoration(labelText: "Username used for login"),
            ),
            SizedBox(height: 10),
            TextField(
              controller: _passwordController,
              decoration: InputDecoration(labelText: "Password used for login"),
            ),
            SizedBox(height: 10),
            TextField(
              controller: _aboutController,
              decoration: InputDecoration(labelText: "About the lecturer"),
            ),
            SizedBox(height: 10),
            ElevatedButton.icon(
              onPressed: () {
                final lec = Lecturer(
                  name: _nameController.text,
                  username: _usernameController.text,
                  password: _passwordController.text,
                  about: _aboutController.text,
                );

                if (lecturer != null) {
                  lec.id = lecturer!.id;
                }

                client.post("/lecturers", data: lec.toJson()).then((res) {
                  if (widget.next != null) widget.next!();
                });
              },
              icon: Icon(lecturer == null ? Icons.add : Icons.edit),
              label: Text("Submit"),
            ),
          ],
        ),
      ),
    );
  }
}
