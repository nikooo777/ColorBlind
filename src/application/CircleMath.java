package application;

import javafx.scene.shape.Circle;

public class CircleMath
{
	// are the circles overlapping?
	static boolean overlap(Circle a, final Circle b)
	{
		return b.getRadius() + a.getRadius() > distance(a, b);
	}

	// returns the distance from one circle to the other
	static double distance(Circle a, final Circle b)
	{
		final double centerXsquared = Math.pow(a.getCenterX() - b.getCenterX(), 2);
		final double centerYsquared = Math.pow(a.getCenterY() - b.getCenterY(), 2);
		return Math.sqrt(centerXsquared + centerYsquared);
	}

	// does the circle intersect another circle?
	public static boolean Intersects(Circle a, Circle b)
	{
		// store the min and max radious of the two circles
		final double minrad = Math.min(a.getRadius(), b.getRadius());
		final double maxrad = Math.max(a.getRadius(), b.getRadius());

		// store the distance between the two circles
		final double distance = distance(a, b);

		// before doing an intense calculation quickly check if they're far apart
		if (minrad + maxrad < distance)
			return false;

		// do intense calculations
		if (distance < Math.abs(maxrad + minrad) && distance(a, b) > Math.abs(maxrad - minrad))
			return true;

		// if nothing is found (they are inside the figure) then return false
		return false;
	}
}