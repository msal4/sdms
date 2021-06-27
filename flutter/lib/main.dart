import 'package:flutter/material.dart';

void main() {
  runApp(App());
}

class App extends StatefulWidget {
  const App({Key? key}) : super(key: key);

  @override
  _AppState createState() => _AppState();
}

class _AppState extends State<App> {
  String _title = 'الصفحه الرئيسيه';

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(
          title: Text(_title),
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
                    )
                  ],
                ),
              ),
              ListTile(
                title: Text('الصفحه الرئيسيه'),
                onTap: () {
                  setState(() {
                    _title = 'الصفحه الرئيسيه';
                  });
                  // Update the state of the app.
                  // ...
                },
              ),
              ListTile(
                title: Text('اهداف القسم'),
                onTap: () {
                  setState(() {
                    _title = 'اهداف القسم';
                  });
                  // Update the state of the app.
                  // ...
                },
              ),
              ListTile(
                title: Text('التدريسيين'),
                onTap: () {
                  setState(() {
                    _title = 'التدريسيين';
                  });
                  // Update the state of the app.
                  // ...
                },
              ),
              ListTile(
                title: Text('المواد الدراسيه'),
                onTap: () {
                  setState(() {
                    _title = 'المواد الدراسيه';
                  });
                  // Update the state of the app.
                  // ...
                },
              ),
              ListTile(
                title: Text('اعلانات'),
                onTap: () {
                  setState(() {
                    _title = 'اعلانات';
                  });
                  // Update the state of the app.
                  // ...
                },
              )
            ],
          ),
        ),
        body: ListView(
          children: [Image.asset("assets/logo.png")],
        ),
      ),
    );
  }
}
