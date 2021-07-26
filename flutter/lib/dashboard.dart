import 'package:flutter/material.dart';
import 'package:sdms/dashboard/pages/announcements.dart';
import 'package:sdms/dashboard/pages/lecturers.dart';
import 'package:sdms/dashboard/pages/login_form.dart';
import 'package:sdms/dashboard/pages/subjects.dart';

void main() {
  runApp(App());
}

class App extends StatelessWidget {
  const App({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      theme: ThemeData.dark(),
      home: LoginFormPage(),
    );
  }
}

class Root extends StatefulWidget {
  const Root({Key? key}) : super(key: key);

  @override
  _RootState createState() => _RootState();
}

class _RootState extends State<Root> with TickerProviderStateMixin {
  int _currentPage = 0;

  final _pages = [LecturersPage(), SubjectsPage(), AnnouncementsPage()];

  get _currentPageWidget => _pages[_currentPage];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(_currentPageWidget.title),
      ),
      drawer: Drawer(
        child: ListView(
          children: [
            Container(
              child: Stack(
                children: [
                  Image.asset(
                    "assets/tech.jpg",
                    fit: BoxFit.cover,
                    // width: double.infinity,
                    // height: double.infinity,
                  ),
                  Positioned(
                    left: 30.0,
                    top: 30.0,
                    child: Container(
                      color: Colors.black.withOpacity(.5),
                      child: Text(
                        'جامعه تكنلوجيا المعلومات والاتصالات',
                        style: TextStyle(
                            color: Colors.white, fontWeight: FontWeight.bold),
                      ),
                    ),
                  ),
                  Positioned(
                    bottom: 0,
                    left: 100,
                    child: Image.asset(
                      "assets/logo.png",
                      fit: BoxFit.cover,
                      width: 100,
                      height: 100,
                    ),
                  ),
                ],
              ),
            ),
            for (int i = 0; i < _pages.length; i++)
              ListTile(
                title: Text((_pages[i] as dynamic).title),
                onTap: () {
                  setState(() {
                    _currentPage = i;
                  });
                  Navigator.pop(context);
                },
              ),
            ListTile(
              title: Text("Logout"),
              onTap: () {
                Navigator.pushReplacement(context,
                    MaterialPageRoute(builder: (context) => LoginFormPage()));
              },
            ),
          ],
        ),
      ),
      body: IndexedStack(children: _pages, index: _currentPage),
    );
  }
}
