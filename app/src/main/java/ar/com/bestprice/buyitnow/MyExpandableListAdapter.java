package ar.com.bestprice.buyitnow;

import android.app.Activity;
import android.util.SparseArray;
import android.view.LayoutInflater;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.ViewGroup;
import android.widget.BaseExpandableListAdapter;
import android.widget.CheckedTextView;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.RelativeLayout;
import android.widget.TextView;
import android.widget.Toast;

import java.util.Map;

import ar.com.bestprice.buyitnow.dto.Item;


public class MyExpandableListAdapter extends BaseExpandableListAdapter {

    private final Map<Integer, PurchasesGroup> groups;
    public LayoutInflater inflater;
    public Activity activity;

    public MyExpandableListAdapter(Activity act, Map<Integer, PurchasesGroup> groups) {
        activity = act;
        this.groups = groups;
        inflater = act.getLayoutInflater();
    }

    @Override
    public Object getChild(int groupPosition, int childPosition) {

        return groups.get(groupPosition).children.get(childPosition);
    }

    @Override
    public long getChildId(int groupPosition, int childPosition) {
        return 0;
    }

    @Override
    public View getChildView(int groupPosition, final int childPosition,
                             boolean isLastChild, View convertView, ViewGroup parent) {

        final Item children = (Item) getChild(groupPosition, childPosition);

        if (convertView == null) {
            convertView = inflater.inflate(R.layout.listrow_details, null);
        }


        TextView text = (TextView) convertView.findViewById(R.id.listrow_item_description);
        text.setText(children.getDescription());

        int icon = Category.MERCADERIA.getIcon();
        if (children.getCategory() != null) {
                icon = children.getCategory().getIcon();
        }

        text.setCompoundDrawablesWithIntrinsicBounds(icon, 0, 0, 0);

        text = (TextView) convertView.findViewById(R.id.item_price);


        text.setText(String.format("%.2f", children.getPrice()) + " $");

        convertView.setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                Toast.makeText(activity, "Category:" + children.getCategory().toString(), Toast.LENGTH_SHORT).show();
            }
        });
        return convertView;
    }

    @Override
    public int getChildrenCount(int groupPosition) {
        return groups.get(groupPosition).children.size();

    }

    @Override
    public Object getGroup(int groupPosition) {
        return groups.get(groupPosition);
    }

    @Override
    public int getGroupCount() {
        return groups.size();
    }

    @Override
    public void onGroupCollapsed(int groupPosition) {
        super.onGroupCollapsed(groupPosition);
    }

    @Override
    public void onGroupExpanded(int groupPosition) {
        super.onGroupExpanded(groupPosition);
    }

    @Override
    public long getGroupId(int groupPosition) {
        return 0;
    }

    @Override
    public View getGroupView(int groupPosition, boolean isExpanded,
                             View convertView, ViewGroup parent) {

        if (convertView == null) {
            convertView = inflater.inflate(R.layout.listrow_group, null);
        }
        //TODO Please try to separate text from numbers by no using tabs
        RelativeLayout relativeLayout = (RelativeLayout) ((LinearLayout)convertView).getChildAt(0);

        CheckedTextView checkedTextView = (CheckedTextView) relativeLayout.getChildAt(0);
        TextView textView = (TextView) relativeLayout.getChildAt(1);
        ImageView image = (ImageView) relativeLayout.getChildAt(2);


        PurchasesGroup purchasesByMonth = (PurchasesGroup) getGroup(groupPosition);
        PurchasesGroup previousPurchasesByMonth = purchasesByMonth;

        if(groupPosition > 0) {
            previousPurchasesByMonth = (PurchasesGroup) getGroup(groupPosition - 1);

           if (purchasesByMonth.getPurchasesTotalPrice() == previousPurchasesByMonth.getPurchasesTotalPrice()) {

                image.setImageResource(R.drawable.icon_minus_24);

           } else if (purchasesByMonth.getPurchasesTotalPrice() > previousPurchasesByMonth.getPurchasesTotalPrice()) {

                image.setImageResource(R.drawable.arrow_up_icon_24);

           } else {
                image.setImageResource(R.drawable.down_icon_24);
           }
        } else {
            image.setImageResource(R.drawable.run_icon_24);
        }
        checkedTextView.setText(purchasesByMonth.getString());
        textView.setText(String.format("$%.2f",purchasesByMonth.getPurchasesTotalPrice()));

        return convertView;
    }

    @Override
    public boolean hasStableIds() {
        return false;
    }

    @Override
    public boolean isChildSelectable(int groupPosition, int childPosition) {
        return false;
    }
}
