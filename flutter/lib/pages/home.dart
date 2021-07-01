import 'package:flutter/material.dart';

import '../const.dart';

class HomePage extends StatefulWidget {
  const HomePage({Key? key}) : super(key: key);

  final title = 'الصفحه الرئيسيه';

  @override
  _HomePageState createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  @override
  Widget build(BuildContext context) {
    return Directionality(
      textDirection: TextDirection.rtl,
      child: ListView(
        padding: const EdgeInsets.all(kDefaultPadding),
        children: [
          Image.asset("assets/logo.png", width: 200, height: 200),
          SizedBox(height: kDefaultPadding * 2),
          Text("قسم تكنولوجیا معلومات الاعمال",
              textAlign: TextAlign.right,
              style: Theme.of(context).textTheme.headline5),
          Divider(),
          Text(
              "یتم اعداد الطالب في قسم تكنولوجیا معلومات الاعمال لیكون مؤھًلا ومواكبا للتطور الحاصل في مجال الادارة وتكنولوجیا المعلومات في آ ٍن واحد.",
              textAlign: TextAlign.right),
          SizedBox(height: kDefaultPadding * 2),
          Text("رؤیة القسم",
              textAlign: TextAlign.right,
              style: Theme.of(context).textTheme.headline5),
          Divider(),
          Text(
              "ان یكون رائداً في مجال التخصص ویتبوأ مكانته بین التخصصات الحدیثة التي تلبي حاجة سوق العمل من ذوي المھارات المتمیزة.",
              textAlign: TextAlign.right),
          SizedBox(height: kDefaultPadding * 2),
          Text("رسالة القسم",
              textAlign: TextAlign.right,
              style: Theme.of(context).textTheme.headline5),
          Divider(),
          Text(
              "یتمتع خریج قسم تكنولوجیا معلومات الاعمال بمھارات تؤھله لان یكون مبرمجا لبورة المھام الإداریة وتوظیف تكنولوجیا المعلومات للعمل الإداري.",
              textAlign: TextAlign.right),
          SizedBox(height: kDefaultPadding * 2),
        ],
      ),
    );
  }
}
