import 'package:flutter/material.dart';
import 'package:sdms/client.dart';
import 'package:sdms/const.dart';
import 'package:sdms/models/subject.dart';
import 'package:url_launcher/url_launcher.dart';

class SubjectsPage extends StatefulWidget {
  const SubjectsPage({Key? key}) : super(key: key);

  final title = 'المواد الدراسيه';

  @override
  _SubjectsPageState createState() => _SubjectsPageState();
}

typedef Semesters = Map<int, List<Subject>>;
typedef Stages = Map<int, Semesters>;

String getYearString(int year) {
  switch (year) {
    case 1:
      return "First";
    case 2:
      return "Second";
    case 3:
      return "Third";
    default:
      return "Fourth";
  }
}

class _SubjectsPageState extends State<SubjectsPage> {
  Stages? _data;

  @override
  void initState() {
    getSubjects().then((subjects) => setState(() {
          _data = subjects.fold({}, (acc, s) {
            if (acc![s.stage] == null) {
              acc[s.stage] = {};
            }
            if (acc[s.stage]![s.semester] == null) {
              acc[s.stage]![s.semester] = [];
            }

            acc[s.stage]![s.semester]!.add(s);

            return acc;
          });
        }));
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    if (_data == null) return Center(child: CircularProgressIndicator());

    return ListView(
      children: [
        for (final stage
            in _data!.keys.toList()..sort((a, b) => a.compareTo(b)))
          ListTile(
            trailing: Icon(Icons.chevron_right),
            title: Text("ISM ${getYearString(stage)} Year"),
            subtitle: Text("${_data![stage]!.keys.length} semesters"),
            leading: Text(stage.toString()),
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (ctx) => SemestersPage(
                    title: "ISM ${getYearString(stage)} Year Semesters",
                    semesters: _data![stage]!,
                  ),
                ),
              );
            },
          )
      ],
    );
  }
}

class SemestersPage extends StatelessWidget {
  const SemestersPage({Key? key, required this.title, required this.semesters})
      : super(key: key);

  final String title;
  final Semesters semesters;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(title),
      ),
      body: ListView(
        children: [
          for (final semester
              in semesters.keys.toList()..sort((a, b) => a.compareTo(b)))
            ListTile(
              trailing: Icon(Icons.chevron_right),
              title: Text("${semester == 1 ? 'First' : 'Second'} Semester"),
              subtitle: Text("${semesters[semester]!.length} classes"),
              leading: Text(semester.toString()),
              onTap: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (ctx) => SemesterSubjectsPage(
                      title:
                          "${semester == 1 ? 'First' : 'Second'} Semester Subjects",
                      subjects: semesters[semester]!,
                    ),
                  ),
                );
              },
            )
        ],
      ),
    );
  }
}

class SemesterSubjectsPage extends StatelessWidget {
  const SemesterSubjectsPage(
      {Key? key, required this.title, required this.subjects})
      : super(key: key);

  final String title;
  final List<Subject> subjects;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(title),
      ),
      body: ListView(
        padding: const EdgeInsets.all(kDefaultPadding),
        children: [
          for (final subject in subjects)
            GestureDetector(
              child: Column(
                children: [
                  Container(
                    decoration: BoxDecoration(
                      color: Colors.grey.shade200,
                      borderRadius: BorderRadius.circular(15),
                    ),
                    padding: const EdgeInsets.all(kDefaultPadding * 2),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(subject.name),
                              SizedBox(height: 10),
                              Container(
                                child: Text(
                                  subject.details ?? "...",
                                  style: TextStyle(color: Colors.grey.shade700),
                                ),
                              ),
                              SizedBox(height: 10),
                              Text(
                                "Lectured by ${subject.lecturer!.name}",
                                style: TextStyle(color: Colors.blue),
                              ),
                            ],
                          ),
                        ),
                        IconButton(
                          tooltip: "Syllabus",
                          icon: Icon(Icons.book),
                          color: Colors.blue,
                          onPressed: () {
                            launch(
                                "http://localhost:5000/storage/pdf/${subject.id}.pdf");
                          },
                        )
                      ],
                    ),
                  ),
                  SizedBox(height: 10),
                ],
              ),
              onTap: () {},
            ),
          SizedBox(height: 50),
        ],
      ),
    );
  }
}
