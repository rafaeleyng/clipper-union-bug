#include "clipper.hpp"
#include <iostream>
#include <vector>

using namespace std;

struct Point {
  double x;
  double y;
  Point(double x, double y) : x(x), y(y) {}
  Point() : x(0), y(0) {}
};

struct Polygon {
  vector<Point> vertices;
  friend ostream &operator<<(ostream &out, const Polygon &polygon) {
    out << "[";
    for (auto &&point : polygon.vertices) {
      out << "{X:" << point.x << " Y:" << point.y << "} ";
    }
    out << "]";
    return out;
  };
  Polygon(size_t n) : vertices(n) {}
  Polygon(vector<Point> points) : vertices{points} {}
  Polygon() : vertices() {}
};

vector<Polygon> getPolygons() {
  Polygon a{{{53, 180},
             {68, 200},
             {44, 199}}};

  Polygon b{{{65, 160},
             {58, 189},
             {30, 190}}};

  Polygon c{{{61, 189},
             {52, 195},
             {48, 187}}};

  vector<Polygon> polygons{a, b, c};
  return polygons;
}

// conversion between clipper types and custom types
const double scaleFactor = 1e9;

ClipperLib::cInt scaleUp(double val) {
  return ClipperLib::cInt(val * scaleFactor);
}

ClipperLib::IntPoint pointToIntPoint(Point p) {
  return ClipperLib::IntPoint(scaleUp(p.x), scaleUp(p.y));
}

ClipperLib::Path polygonToPath(Polygon p) {
  ClipperLib::Path newPath(p.vertices.size());
  for (size_t i = 0; i < p.vertices.size(); i++) {
    newPath[i] = pointToIntPoint(p.vertices[i]);
  }
  return newPath;
}

ClipperLib::Paths polygonsToPaths(vector<Polygon> polygons) {
  ClipperLib::Paths newPaths(polygons.size());
  for (size_t i = 0; i < polygons.size(); i++) {
    newPaths[i] = polygonToPath(polygons[i]);
  }
  return newPaths;
}

double scaleDown(ClipperLib::cInt val) { return double(val) / scaleFactor; }

Point intPointToPoint(ClipperLib::IntPoint p) {
  return Point(scaleDown(p.X), scaleDown(p.Y));
}

Polygon pathToPolygon(ClipperLib::Path path) {
  Polygon newPoly(path.size());
  for (size_t i = 0; i < path.size(); i++) {
    newPoly.vertices[i] = intPointToPoint(path[i]);
  }
  return newPoly;
}

vector<Polygon> pathsToPolygons(ClipperLib::Paths paths) {
  vector<Polygon> newPolygons(paths.size());
  for (size_t i = 0; i < paths.size(); i++) {
    newPolygons[i] = pathToPolygon(paths[i]);
  }
  return newPolygons;
}

// run with: g++ -Wall -std=c++17 -o main main.cpp clipper.cpp
int main() {
  ClipperLib::Clipper clip;
  auto polygons = getPolygons();
  auto paths = polygonsToPaths(polygons);

//  ClipperLib::SimplifyPolygons(paths, ClipperLib::pftNonZero);
  clip.AddPaths(paths, ClipperLib::ptSubject, true);

  ClipperLib::Paths combinedPaths;
  const auto ok = clip.Execute(ClipperLib::ctUnion, combinedPaths,
                               ClipperLib::pftNonZero, ClipperLib::pftNonZero);

  if (!ok || combinedPaths.empty() || combinedPaths.size() > 1) {
    cout << "Path union failed! expected one path but got: "
         << combinedPaths.size() << endl;
  }

  auto combinedPolys = pathsToPolygons(combinedPaths);
  for (auto &&p : combinedPolys) {
    cout << "[" << p << "]" << endl;
  }
}

// Output:
// {
// [X: 59.573,Y: 195.326]
// [X: 102.667,Y: 213.998]
// [X: 180.578,Y: 238.897]
// [X: 178.244,Y: 266.02]
// [X: 59.9828,Y: 272.401]
// [X: 58.691,Y: 287.413]
// [X: 59.3193,Y: 288.502]
// [X: 61.025,Y: 288.705]
// [X: 86.698,Y: 295.342]
// [X: 89.294,Y: 312.188]
// [X: 86.96,Y: 339.311]
// [X: 73.734,Y: 345.027]
// [X: 58.7773,Y: 343.241]
// [X: 58.485,Y: 346.637]
// [X: 25.673,Y: 346.637]
// [X: 1,Y: 340]
// [X: 1,Y: 323.154]
// [X: 5.65111,Y: 312.836]
// [X: 1,Y: 311.585]
// [X: 1,Y: 294.739]
// [X: 5.65111,Y: 284.421]
// [X: 1,Y: 283.17]
// [X: 1,Y: 266.324]
// [X: 9.62435,Y: 247.192]
// [X: 1,Y: 238.116]
// [X: 1,Y: 221.27]
// [X: 28.97,Y: 198.966]
// [X: 42.196,Y: 193.25]
// }