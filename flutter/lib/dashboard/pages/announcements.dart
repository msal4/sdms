import 'package:flutter/material.dart';
import 'package:sdms/client.dart';
import 'package:sdms/const.dart';
import 'package:sdms/dashboard/pages/announcement_form.dart';
import 'package:sdms/models/announcement.dart';

class AnnouncementsPage extends StatefulWidget {
  const AnnouncementsPage({Key? key}) : super(key: key);

  final title = "Manage Announcements";

  @override
  _AnnouncementsPageState createState() => _AnnouncementsPageState();
}

class _AnnouncementsPageState extends State<AnnouncementsPage> {
  List<Announcement> _data = [];

  @override
  void initState() {
    getAnnouncements().then((value) => setState(() {
          _data = value;
        }));
    super.initState();
  }

  refetch() => getAnnouncements().then(
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
                    builder: (ctx) => AnnouncementFormPage(next: refetch),
                  ),
                );
              },
              icon: Icon(Icons.add),
              label: Text("Add Announcement")),
          for (final item in _data)
            ListTile(
              title: Text(item.title),
              leading: IconButton(
                  onPressed: () {
                    client
                        .delete("/announcements/${item.id}")
                        .then((value) => refetch());
                  },
                  icon: Icon(Icons.delete, color: Colors.red)),
            )
        ],
      ),
    );
  }
}
