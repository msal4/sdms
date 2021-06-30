import 'package:flutter/material.dart';
import 'package:sdms/client.dart';
import 'package:sdms/const.dart';
import 'package:sdms/models/announcement.dart';

class AnnouncementsPage extends StatefulWidget {
  const AnnouncementsPage({Key? key}) : super(key: key);

  final title = 'اعلانات القسم';

  @override
  _AnnouncementsPageState createState() => _AnnouncementsPageState();
}

class _AnnouncementsPageState extends State<AnnouncementsPage> {
  List<Announcement>? _data;

  @override
  void initState() {
    getAnnouncements().then((announcements) => setState(() {
          _data = announcements;
        }));
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    if (_data == null) return Center(child: CircularProgressIndicator());

    return ListView(
      padding: const EdgeInsets.all(kDefaultPadding),
      children: [
        for (final announcement in _data!)
          Column(
            children: [
              Container(
                clipBehavior: Clip.antiAlias,
                decoration: BoxDecoration(
                  color: Colors.grey.shade200,
                  borderRadius: BorderRadius.circular(15),
                ),
                child: Column(
                  children: [
                    Container(
                      padding: const EdgeInsets.all(kDefaultPadding * 2),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(announcement.title),
                                SizedBox(height: 10),
                                Container(
                                  child: Text(
                                    announcement.details ?? "...",
                                    style:
                                        TextStyle(color: Colors.grey.shade700),
                                  ),
                                ),
                              ],
                            ),
                          ),
                          Icon(Icons.notifications),
                        ],
                      ),
                    ),
                    Image.network(
                      "http://localhost:5000/storage/images/6.png",
                      height: 180,
                      fit: BoxFit.cover,
                      width: double.infinity,
                    )
                  ],
                ),
              ),
              SizedBox(height: 10),
            ],
          ),
      ],
    );
  }
}
