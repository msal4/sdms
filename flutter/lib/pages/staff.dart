import 'package:flutter/material.dart';
import 'package:sdms/client.dart';
import 'package:sdms/models/lecturer.dart';
import 'package:sdms/pages/subjects.dart';

class LecturersPage extends StatefulWidget {
  const LecturersPage({Key? key}) : super(key: key);

  final title = 'الكادر التدريسي';

  @override
  _LecturersPageState createState() => _LecturersPageState();
}

class _LecturersPageState extends State<LecturersPage> {
  List<Lecturer>? _data;

  @override
  void initState() {
    getLecturers().then((lecturers) => setState(() {
          _data = lecturers;
        }));
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    if (_data == null) return Center(child: CircularProgressIndicator());

    return ListView(
      children: [
        for (final lecturer in _data!)
          ListTile(
            trailing: Icon(Icons.chevron_right),
            title: Text(lecturer.name),
            subtitle: Text("@${lecturer.username}"),
            leading: Icon(Icons.person),
            onTap: () {
              getSubjectsByLecturerID(lecturer.id!).then(
                (subjects) {
                  return Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (ctx) => SemesterSubjectsPage(
                          title: "${lecturer.name}'s Subjects",
                          subjects: subjects),
                    ),
                  );
                },
              );
            },
          )
      ],
    );
  }
}
