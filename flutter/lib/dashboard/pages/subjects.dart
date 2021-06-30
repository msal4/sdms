import 'package:flutter/material.dart';
import 'package:sdms/client.dart';
import 'package:sdms/const.dart';
import 'package:sdms/dashboard/pages/subject_form.dart';
import 'package:sdms/models/subject.dart';

class SubjectsPage extends StatefulWidget {
  const SubjectsPage({Key? key}) : super(key: key);

  final title = "Manage Subjects";

  @override
  _SubjectsPageState createState() => _SubjectsPageState();
}

class _SubjectsPageState extends State<SubjectsPage> {
  List<Subject> _data = [];

  @override
  void initState() {
    getSubjects().then((value) => setState(() {
          _data = value;
        }));
    super.initState();
  }

  refetch() => getSubjects().then(
        (value) => setState(() {
          _data = value;
        }),
      );

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(kDefaultPadding),
      child: ListView(
        children: [
          ElevatedButton.icon(
              onPressed: () {
                Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (ctx) => SubjectFormPage(next: refetch),
                  ),
                );
              },
              icon: Icon(Icons.add),
              label: Text("Add Subject")),
          for (final item in _data)
            ListTile(
              title: Text(item.name),
              subtitle:
                  Text("semester: ${item.semester}, stage: ${item.stage}"),
              leading: IconButton(
                  onPressed: () {
                    client
                        .delete("/subjects/${item.id}")
                        .then((value) => refetch());
                  },
                  icon: Icon(Icons.delete, color: Colors.red)),
            )
        ],
      ),
    );
  }
}
