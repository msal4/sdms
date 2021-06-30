import 'package:flutter/material.dart';
import 'package:sdms/client.dart';
import 'package:sdms/const.dart';
import 'package:sdms/dashboard/pages/lecturer_form.dart';
import 'package:sdms/models/lecturer.dart';

class LecturersPage extends StatefulWidget {
  const LecturersPage({Key? key}) : super(key: key);

  final title = "Manage Lecturers";

  @override
  _LecturersPageState createState() => _LecturersPageState();
}

class _LecturersPageState extends State<LecturersPage> {
  List<Lecturer> _data = [];

  @override
  void initState() {
    getLecturers().then((value) => setState(() {
          _data = value;
        }));
    super.initState();
  }

  refetch() => getLecturers().then(
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
                    builder: (ctx) => LecturerFormPage(next: refetch),
                  ),
                );
              },
              icon: Icon(Icons.add),
              label: Text("Add Lecturer")),
          for (final item in _data)
            ListTile(
              title: Text(item.name),
              subtitle: Text("@" + item.username),
              leading: IconButton(
                  onPressed: () {
                    client
                        .delete("/lecturers/${item.id}")
                        .then((value) => refetch());
                  },
                  icon: Icon(Icons.delete, color: Colors.red)),
              trailing: IconButton(
                  onPressed: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (ctx) => LecturerFormPage(
                          next: () {
                            refetch();
                            Navigator.pop(context);
                          },
                          lecturer: item,
                        ),
                      ),
                    );
                  },
                  icon: Icon(Icons.edit)),
            )
        ],
      ),
    );
  }
}
