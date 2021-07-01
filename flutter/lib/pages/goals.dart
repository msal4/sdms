import 'package:flutter/material.dart';
import 'package:sdms/const.dart';

class GoalsPage extends StatefulWidget {
  const GoalsPage({Key? key}) : super(key: key);

  final title = 'اهداف القسم';

  @override
  _GoalsPageState createState() => _GoalsPageState();
}

class _GoalsPageState extends State<GoalsPage> {
  @override
  Widget build(BuildContext context) {
    return Directionality(
      textDirection: TextDirection.rtl,
      child: ListView(
        padding: const EdgeInsets.all(kDefaultPadding),
        children: [
          Row(
            children: [
              Image.asset("assets/logo.png", width: 70, height: 70),
              SizedBox(width: kDefaultPadding),
              Text("اھداف قسم تكنولوجیا معلومات الاعمال",
                  textAlign: TextAlign.center,
                  style: Theme.of(context).textTheme.headline5),
            ],
          ),
          Divider(),
          Text(
              "یتم اعداد الطالب في قسم تكنولوجیا معلومات الاعمال لیكون مؤھًلا ومواكبا للتطور الحاصل في مجال الادارة وتكنولوجیا المعلومات في آٍن واحد.",
              textAlign: TextAlign.right),
          SizedBox(height: kDefaultPadding * 2),
          ListTile(
            leading: Icon(Icons.home_work_outlined),
            title: Text(
              "توفیر البنى التحتیة لبیئة تكنولوجیا المعلومات",
              textAlign: TextAlign.right,
            ),
          ),
          SizedBox(height: kDefaultPadding * 2),
          ListTile(
            leading: Icon(Icons.security_outlined),
            title: Text(
              "تحسین عملیة إدارة وامن تكنولوجیا المعلومات",
              textAlign: TextAlign.right,
            ),
          ),
          SizedBox(height: kDefaultPadding * 2),
          ListTile(
            leading: Icon(Icons.leaderboard_outlined),
            title: Text(
              "رفد سوق العمل بما یتطلب من كوادر إداریة واستراتیجیة",
              textAlign: TextAlign.right,
            ),
          ),
          SizedBox(height: kDefaultPadding * 2),
          ListTile(
            leading: Icon(Icons.school_outlined),
            title: Text(
              "تخریج كوادر لھم القابلیة على التعلم ومواكبه التطور ومنافسة اقرانھم عالمیا",
              textAlign: TextAlign.right,
            ),
          ),
        ],
      ),
    );
  }
}
